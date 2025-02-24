<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;

class TestController 
{
    #[Route("/test_plain")]
    public function test_handler(): Response {
        $resp = new Response("Hello world!");
        $resp->headers->add([
            "Content-Type" => "text/plain"
        ]);
        return $resp;
    }
}

?>