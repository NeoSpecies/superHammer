<?php
namespace Develop\Reader;

use Exception;
class UnixSocketReader {
    protected $pool = [];
    protected $socketFile;
    protected $poolSize;
    private $container;

    public function __construct($container) {
        $this->socketFile = $container->get('socketMainFile');
        $this->poolSize = 5;

        for ($i = 0; $i < $this->poolSize; $i++) {
            $this->pool[$i] = $this->createConnection();
        }
    }

    protected function createConnection() {
        $socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
        if (!$socket) {
            throw new \RuntimeException("Unable to create socket: " . socket_strerror(socket_last_error()));
        }
        if (!socket_connect($socket, $this->socketFile)) {
            socket_close($socket);
            throw new \RuntimeException("Unable to connect: " . socket_strerror(socket_last_error($socket)));
        }
        return $socket;
    }

    public function getConnection() {
        foreach ($this->pool as $key => $socket) {
            if (is_resource($socket) && !socket_get_option($socket, SOL_SOCKET, SO_ERROR)) {
                return $socket;
            } else {
                // Close the invalid socket connection
                if (is_resource($socket)) {
                    socket_close($socket);
                }
                // Try to recreate the socket if it is not valid anymore
                $this->pool[$key] = $this->createConnection();
                return $this->pool[$key];
            }
        }
        // If all sockets are busy, create a new one
        $newSocket = $this->createConnection();
        $this->pool[] = $newSocket;
        return $newSocket;
    }

    public function release($socket) {
        // For now, do nothing. Add back to pool or recreate if needed.
    }

    public function sendAndReceive($message) {
        try {
            $socket = $this->getConnection();
            if (socket_write($socket, $message, strlen($message)) === false) {
                throw new \RuntimeException("Failed to write to socket: " . socket_strerror(socket_last_error($socket)));
            }
            $response = socket_read($socket, 2048);
            if ($response === false) {
                throw new \RuntimeException("Failed to read from socket: " . socket_strerror(socket_last_error($socket)));
            }
            $this->release($socket);
            return $response;
        } catch (\Exception $e) {
            // Handle exception by logging or rethrowing
            error_log($e->getMessage());
            throw $e;
        }
    }

    public function __destruct() {
        foreach ($this->pool as $socket) {
            if (is_resource($socket)) {
                socket_close($socket);
            }
        }
    }
}
