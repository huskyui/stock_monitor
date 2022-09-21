package com.stock.stockdata.service.impl;

import com.stock.stockdata.model.StockInfo;
import com.stock.stockdata.service.StockInfoService;
import com.stock.stockdata.util.DateUtil;
import com.stock.stockdata.util.HttpClientUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.util.Date;

@Service
public class StockInfoServiceImpl implements StockInfoService {
    @Override
    public StockInfo getStockRealTimeInfo(String stockNum) {

        String url = "http://qt.gtimg.cn/q=s_" + stockNum;

        String response = HttpClientUtil.get(url);
        if (StringUtils.isEmpty(response)) {
            return null;
        }

        int index = response.indexOf("=");
        String needInfo = response.substring(index + 2, response.length() - 3);
        String[] stockInfoArr = needInfo.split("~");
        StockInfo stockInfo = new StockInfo();
        stockInfo.setStockNum(stockNum);
        stockInfo.setName(stockInfoArr[1]);
        stockInfo.setPrice(new BigDecimal(stockInfoArr[3]));
        stockInfo.setPriceDiff(new BigDecimal(stockInfoArr[4]));
        stockInfo.setDiffPercent(new BigDecimal(stockInfoArr[5]));
        stockInfo.setDate(DateUtil.getDayStartDate(new Date()));
        return stockInfo;
    }


    public static void main(String[] args) {
        StockInfoServiceImpl stockInfoService = new StockInfoServiceImpl();
        System.out.println(stockInfoService.getStockRealTimeInfo("sz002497"));
    }

}
