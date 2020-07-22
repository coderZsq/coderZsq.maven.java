package com.sq.productserver.web.feign;

import com.sq.productapi.domain.Product;
import com.sq.productapi.feign.ProductFeignApi;
import com.sq.productserver.service.IProductService;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.RestController;

import java.util.concurrent.TimeUnit;

@RestController
public class ProductFeignClient implements ProductFeignApi {
    @Autowired
    private IProductService productService;

    @Value("${server.port}")
    private String port;

    @Override
    public Product find(Long id) {
        System.out.println("调用商品服务feign接口");
        // 业务方法执行, 重试策略 --> 业务方法会执行多次 --> 接口一定要做幂等性
        Product product = productService.get(id);
        try {
            TimeUnit.SECONDS.sleep(5);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        Product result = new Product();
        BeanUtils.copyProperties(product, result);
        result.setName(result.getName() + ", data from " + port);
        return result;
    }
}
