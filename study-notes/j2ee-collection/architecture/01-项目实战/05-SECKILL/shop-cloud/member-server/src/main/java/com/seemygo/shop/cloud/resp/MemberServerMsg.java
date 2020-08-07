package com.seemygo.shop.cloud.resp;

public class MemberServerMsg extends CodeMsg {

    public static final MemberServerMsg PARAM_ERROR = new MemberServerMsg("500400", "参数错误：{0}");
    public static final MemberServerMsg USERNAME_OR_PASSWORD_ERROR = new MemberServerMsg("500401", "用户名或密码错误");

    private MemberServerMsg(String code, String msg) {
        super(code, msg);
    }
}
