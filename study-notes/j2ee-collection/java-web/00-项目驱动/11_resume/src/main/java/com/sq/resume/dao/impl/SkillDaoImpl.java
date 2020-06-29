package com.sq.resume.dao.impl;

import com.sq.resume.bean.Skill;
import com.sq.resume.dao.SkillDao;
import com.sq.resume.util.Dbs;
import org.springframework.jdbc.core.BeanPropertyRowMapper;

import java.util.ArrayList;
import java.util.List;

public class SkillDaoImpl extends BaseDaoImpl<Skill> implements SkillDao {
    @Override
    protected String table() {
        return "skill";
    }

    @Override
    public boolean save(Skill bean) {
        Integer id = bean.getId();
        List<Object> args = new ArrayList<>();
        args.add(bean.getName());
        args.add(bean.getLevel());
        String sql;
        if (id == null || id < 1) { // 添加
            sql = "INSERT INTO skill(name, level) VALUES(?, ?)";
        } else {
            sql = "UPDATE skill SET name = ?, level = ? WHERE id = ?";
            args.add(id);
        }
        return Dbs.getTpl().update(sql, args.toArray()) > 0;
    }

    @Override
    public Skill get(Integer id) {
        String sql = "SELECT id, created_time, name, level FROM skill WHERE id = ?";
        return Dbs.getTpl().queryForObject(sql, new BeanPropertyRowMapper<>(Skill.class), id);
    }

    /**
     * 获取多个对象
     */
    @Override
    public List<Skill> list() {
        String sql = "SELECT id, created_time, name, level FROM skill";
        return Dbs.getTpl().query(sql, new BeanPropertyRowMapper<>(Skill.class));
    }
}
