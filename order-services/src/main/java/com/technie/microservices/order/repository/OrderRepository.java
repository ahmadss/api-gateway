package com.technie.microservices.order.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import com.technie.microservices.order.model.Order;

public interface OrderRepository extends JpaRepository<Order, Long>{

}
