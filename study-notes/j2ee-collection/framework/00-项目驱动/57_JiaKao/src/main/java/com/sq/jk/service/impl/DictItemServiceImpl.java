package com.sq.jk.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.sq.jk.common.enhance.MpPage;
import com.sq.jk.common.enhance.MpQueryWrapper;
import com.sq.jk.mapper.DictItemMapper;
import com.sq.jk.pojo.po.DictItem;
import com.sq.jk.pojo.query.DictItemQuery;
import com.sq.jk.service.DictItemService;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Transactional
public class DictItemServiceImpl extends ServiceImpl<DictItemMapper, DictItem> implements DictItemService {

    @Override
    public void list(DictItemQuery query) {
        // 查询条件
        MpQueryWrapper<DictItem> wrapper = new MpQueryWrapper<>();
        wrapper.like(query.getKeyword(), DictItem::getName, DictItem::getValue);
        Integer typeId = query.getTypeId();
        if (typeId != null && typeId > 0) {
            wrapper.eq(DictItem::getTypeId, typeId);
        }

        // 排序
        wrapper.orderByDesc(DictItem::getId);

        // 查询
        baseMapper.selectPage(new MpPage<>(query), wrapper).updateQuery();
    }
}