package com.technie.microservices.order.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.swagger.v3.oas.models.ExternalDocumentation;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;

@Configuration
public class OpenAPIConfig {

	@Bean
	public OpenAPI productServiceApi() {
		return new OpenAPI()
				.info(new Info().title("Order Service API")
						.description("This is the rest Api for Order Service")
						.version("v0.0.1")
						.license(new License().name("Apache 2.0")))
				.externalDocs(new ExternalDocumentation().description("You can refer to the order service wiki documentation")
						.url("https://product-service-dummy-url.com/docs"));
	}
}
