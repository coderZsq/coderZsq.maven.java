package com.seemygo.shop.cloud.mapper;

import com.seemygo.shop.cloud.domain.OrderInfo;
import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;

@Mapper
public interface OrderInfoMapper {
    @Insert("insert into t_order_info(order_no,user_id,good_id,delivery_addr_id,good_name,good_img,good_count,good_price,seckill_price,status,create_date,pay_date) values(#{orderNo},#{userId},#{goodId},#{deliveryAddrId},#{goodName},#{goodImg},#{goodCount},#{goodPrice},#{seckillPrice},#{status},#{createDate},#{payDate})")
    void insert(OrderInfo orderInfo);

    @Select("select * from t_order_info where order_no = #{orderNo} AND user_id = #{userId}")
    OrderInfo selectByPrimaryKey(@Param("orderNo") String orderNo, @Param("userId") Long userId);
}
