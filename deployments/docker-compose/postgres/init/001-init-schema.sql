-- Create custom schema
CREATE SCHEMA IF NOT EXISTS rockets_message_processor AUTHORIZATION rockets_message_processor_role;

-- Set default search_path for the user
ALTER ROLE rockets_message_processor_role SET search_path TO rockets_message_processor;