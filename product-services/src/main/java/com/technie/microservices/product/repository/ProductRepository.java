package com.technie.microservices.product.repository;

import org.springframework.data.mongodb.repository.MongoRepository;

import com.technie.microservices.product.model.Product;

public interface ProductRepository extends MongoRepository<Product, String> {
}
