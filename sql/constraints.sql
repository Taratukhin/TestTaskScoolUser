DROP PROCEDURE IF EXISTS PROC_DROP_CONSTRAINT;

CREATE PROCEDURE PROC_DROP_CONSTRAINT(IN tableName VARCHAR(64), IN constraintName VARCHAR(64), IN constraintType VARCHAR(64))
    BEGIN
        IF EXISTS(
            SELECT * FROM information_schema.table_constraints
            WHERE 
                table_schema    = DATABASE()     AND
                table_name      = tableName      AND
                constraint_name = constraintName)
        THEN
            SET @query = CONCAT('ALTER TABLE ', tableName, ' DROP ', constraintType, ' ', constraintName, ';');
            PREPARE stmt FROM @query; 
            EXECUTE stmt; 
            DEALLOCATE PREPARE stmt; 
        END IF; 
    END;

CALL PROC_DROP_CONSTRAINT('user', 'UC_id','CONSTRAINT');
ALTER TABLE user ADD CONSTRAINT UC_id UNIQUE(username); -- in this task, the username must be unique, as it identifies the user

CALL PROC_DROP_CONSTRAINT('user_profile', 'UC_user_id','CONSTRAINT');
ALTER TABLE user_profile ADD CONSTRAINT UC_user_id   UNIQUE(user_id); -- in this task, the user only attends one school
CALL PROC_DROP_CONSTRAINT('user_profile', 'UC_user_attr','CONSTRAINT');
ALTER TABLE user_profile ADD CONSTRAINT UC_user_attr UNIQUE(first_name,last_name, phone, address,city);
CALL PROC_DROP_CONSTRAINT('user_profile', 'FK_user_id','FOREIGN KEY');
ALTER TABLE user_profile ADD CONSTRAINT FK_user_id FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE;

CALL PROC_DROP_CONSTRAINT('user_data', 'UC_user_id','CONSTRAINT');
ALTER TABLE user_data ADD CONSTRAINT UC_user_id   UNIQUE(user_id); -- in this task, one user has only one first name, last name, address and city
CALL PROC_DROP_CONSTRAINT('user_data', 'FK_user_data_id','FOREIGN KEY');  -- This task cannot have two different users with the same attributes
ALTER TABLE user_data ADD CONSTRAINT FK_user_data_id   FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE;
