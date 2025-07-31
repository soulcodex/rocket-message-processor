-- Create rocket data table
CREATE TABLE "rockets_message_processor.rocket" (
    id           VARCHAR(50) PRIMARY KEY,
    rocket_type  VARCHAR(100) NOT NULL,
    launch_speed INTEGER      NOT NULL CHECK (launch_speed > 0),
    mission      TEXT         NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL,
    updated_at   TIMESTAMPTZ  NOT NULL
);
