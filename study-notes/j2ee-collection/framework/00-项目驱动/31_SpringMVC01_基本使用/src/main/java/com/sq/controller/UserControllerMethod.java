package com.sq.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

@Controller
@RequestMapping("/user")
public class UserControllerMethod {
    // /user/add
    // @RequestMapping(value = "/add", method = RequestMethod.POST)
    @PostMapping("/add")
    @ResponseBody
    public String add() {
        return "UserController - Add Success!";
    }

    // /user/get
    // @RequestMapping(value = "/get", method = RequestMethod.GET)
    @GetMapping({"/get", "/list"})
    @ResponseBody
    public String get() {
        return "UserController - Get Success!";
    }
}
