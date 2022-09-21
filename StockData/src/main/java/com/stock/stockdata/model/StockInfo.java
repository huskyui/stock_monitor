package com.stock.stockdata.model;

import lombok.Data;

import java.math.BigDecimal;
import java.util.Date;

@Data
public class StockInfo {
    private String name;

    private String stockNum;

    private BigDecimal price;

    private BigDecimal diffPercent;

    private BigDecimal priceDiff;

    private Date date;
}
