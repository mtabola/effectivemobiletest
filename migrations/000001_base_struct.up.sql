CREATE TABLE IF NOT EXISTS Regions (
    RegionId serial PRIMARY KEY,
    RegionCode varchar(3) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Genders (
    GenderId serial PRIMARY KEY,
    GenderName varchar(16) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Users (
    UserId serial PRIMARY KEY,
    FName varchar(64) NOT NULL,
    SName varchar(64) NOT NULL,
    PName varchar(64),
    GenderId integer NOT NULL,
    Age integer NOT NULL check (Age > 0),
    RegionId integer NOT NULL,
    CONSTRAINT FK_Region FOREIGN KEY (RegionId) REFERENCES Regions(RegionId),
    CONSTRAINT FK_Gender FOREIGN KEY (GenderId) REFERENCES Genders(GenderId)
);