CREATE TABLE Users (
    UserID STRING(MAX) NOT NULL,
    Name STRING(MAX) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true ),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true ),
) PRIMARY KEY (UserID);
