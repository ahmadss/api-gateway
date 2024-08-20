package com.technie.microservices.inventory;

import org.springframework.boot.SpringApplication;

public class TestInventoryServicesApplication {

	public static void main(String[] args) {
		SpringApplication.from(InventoryServicesApplication::main).with(TestcontainersConfiguration.class).run(args);
	}

}
