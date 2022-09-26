package com.stock.stockdata.api;

import com.stock.stockdata.db.Db;
import com.stock.stockdata.model.StockInfo;
import com.stock.stockdata.service.StockInfoService;
import com.stock.stockdata.util.DateUtil;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Resource;
import java.io.UnsupportedEncodingException;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api")
public class StockApi {
    @Resource
    private StockInfoService stockInfoService;

    @Resource
    private Db db;

    @RequestMapping("/getStockInfo")
    public String getStockInfo(@RequestParam("stockNum") String stockNum) throws SQLException, UnsupportedEncodingException {
        StockInfo stockInfo = stockInfoService.getStockRealTimeInfo(stockNum);
        Map<String, Object> data = new HashMap<>();
        data.put("name", stockInfo.getName());
        data.put("stock_num", stockInfo.getStockNum());
        data.put("price", stockInfo.getPrice());
        data.put("diff_percent", stockInfo.getDiffPercent());
        data.put("price_diff", stockInfo.getPriceDiff());
        data.put("now_time", DateUtil.date2String(stockInfo.getDate()));


        List<Map<String, Object>> datas = new ArrayList<>();
        datas.add(data);
        db.insertAll("t_stock", datas);
        return "success";


    }

}
