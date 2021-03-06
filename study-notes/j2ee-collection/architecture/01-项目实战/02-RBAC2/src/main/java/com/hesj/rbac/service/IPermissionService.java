package com.hesj.rbac.service;

import com.hesj.rbac.domain.Permission;
import com.hesj.rbac.qo.PageResult;
import com.hesj.rbac.qo.QueryObject;

import java.util.List;

public interface IPermissionService {
    //添加和更新
    public void saveOrUpdate(Permission permission);

    //删除
    void delete(Long id);


    //根据id查
    Permission getById(Long id);

    //查所有
    List<Permission> getAllList();

    PageResult<Permission> selectByQuery(QueryObject qo);

    void reload();

    List<Permission> getAllPermissionByRoleId(Long roleId);

    List<String> selectAllExpressionsByEmployeeId(Long employeeId);
}
