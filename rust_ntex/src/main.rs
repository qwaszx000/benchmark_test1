use ntex::web;

#[web::get("/test_plain")]
async fn test_handler() -> impl web::Responder {
    web::HttpResponse::Ok().body("Hello world!")
}

#[ntex::main]
async fn main() -> std::io::Result<()> {
    web::HttpServer::new(|| {
        web::App::new()
            .service(test_handler)
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
