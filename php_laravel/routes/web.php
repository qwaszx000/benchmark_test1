<?php

use Illuminate\Support\Facades\Route;
use Illuminate\Support\Facades\Response;

Route::get('/test_plain', function () {
    $resp = Response::make("Hello world!");
    $resp->header("Content-Type", "text/plain");
    return $resp;
});