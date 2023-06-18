CREATE TABLE IF NOT EXISTS ALL_REQUESTS
(
    ID            BIGINT PRIMARY KEY NOT NULL, /*TIMESTAMP IN NANO SECONDS*/
    TUNNEL_ID     TEXT, /* GOLANG UUID */
    IS_REPLAY     BOOLEAN,/*IS THIS A REPLAYED REQUEST OR NOT*/
    REQUEST_DATA  BLOB,
    RESPONSE_DATA BLOB,
    LOCAL_PORT    INT,
    METADATA      TEXT DEFAULT '{}'
);
