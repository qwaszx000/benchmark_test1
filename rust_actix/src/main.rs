use actix_web;


#[actix_web::get("/test_plain")]
async fn test_handler() -> impl actix_web::Responder {
    actix_web::HttpResponse::Ok().body("Hello world!")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    actix_web::HttpServer::new(|| {
        actix_web::App::new()
            .service(test_handler)
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
