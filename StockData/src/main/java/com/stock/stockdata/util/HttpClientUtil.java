package com.stock.stockdata.util;

import org.apache.http.HttpEntity;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;

import java.nio.charset.StandardCharsets;

public class HttpClientUtil {
    private static final CloseableHttpClient httpClient;


    static {
        httpClient = HttpClients.createDefault();
    }

    public static String get(String url) {
        HttpGet httpGet = new HttpGet(url);
        String responseBody = null;
        try {
            responseBody = httpClient.execute(httpGet, response -> {
                int statusCode = response.getStatusLine().getStatusCode();
                if (statusCode >= 200 && statusCode < 300) {
                    HttpEntity entity = response.getEntity();
                    return entity != null ? EntityUtils.toString(entity, StandardCharsets.UTF_8) : null;
                } else {
                    throw new ClientProtocolException("Unexpected response status: " + statusCode);
                }
            });
        } catch (Exception e) {

        }
        return responseBody;

    }


}
