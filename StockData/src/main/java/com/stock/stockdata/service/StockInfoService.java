package com.stock.stockdata.service;

import com.stock.stockdata.model.StockInfo;


public interface StockInfoService {
    StockInfo getStockRealTimeInfo(String stockNum);

}
