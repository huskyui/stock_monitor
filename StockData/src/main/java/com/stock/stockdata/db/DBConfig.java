package com.stock.stockdata.db;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import ru.yandex.clickhouse.ClickHouseDataSource;
import ru.yandex.clickhouse.settings.ClickHouseProperties;

import java.sql.Connection;
import java.sql.SQLException;

@Configuration
public class DBConfig {
    @Value("${clickhouse.url}")
    private String address;
    @Value("${clickhouse.username}")
    private String username;
    @Value("${clickhouse.db}")
    private String db;


    public Connection getConnection() {
        ClickHouseProperties properties = new ClickHouseProperties();
        properties.setUser(username);
        properties.setDatabase(db);
        properties.setSocketTimeout(60000);
        ClickHouseDataSource clickHouseDataSource = new ClickHouseDataSource(address, properties);

        try {
            return clickHouseDataSource.getConnection();
        } catch (SQLException e) {
            e.printStackTrace();
            return null;
        }
    }

}
