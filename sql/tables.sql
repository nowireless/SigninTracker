--  docker run --name signin-postgres -p5432:5432 -e POSTGRES_USER=signin -e POSTGRES_PASSWORD=foobar -d postgres

DROP TYPE IF EXISTS COMPETITION;
CREATE TYPE COMPETITION as ENUM('FRC', 'FTC', 'FLL');

DROP TYPE IF EXISTS PARENTRELATION;
CREATE TYPE PARENTRELATION as ENUM('Father', 'Mother', 'Guardian');

DROP TABLE IF EXISTS people;
CREATE TABLE people (
  PersonID SERIAL PRIMARY KEY, --INTEGER PRIMARY KEY,
  CheckInID TEXT UNIQUE NOT NULL,
  FirstName TEXT NOT NULL,
  LastName  TEXT NOT NULL,
  Email     TEXT,
  Phone     TEXT,

  SchoolEmail TEXT,
  SchoolID    TEXT
);

DROP TABLE IF EXISTS meetings;
CREATE TABLE meetings(
  MeetingID SERIAL PRIMARY KEY, --INTEGER PRIMARY KEY,

  Date      DATE NOT NULL,
  StartTime TIME NOT NULL,
  EndTime   TIME NOT NULL,
  Location  TEXT NOT NULL, -- TOOD make location its own table?
  Kind      TEXT NOT NULL
);

DROP TABLE IF EXISTS teams;
CREATE TABLE teams (
  TeamID SERIAL PRIMARY KEY,

  Compeition  COMPETITION NOT NULL, -- FIX typo
  Number      INTEGER NOT NULL,
  Name        TEXT NOT NULL,

  UNIQUE (Compeition, Number)
);

DROP TABLE IF EXISTS team_meetings;
CREATE TABLE team_meetings(
  TeamID INTEGER NOT NULL REFERENCES teams(TeamID),
  MeetingID INTEGER NOT NULL REFERENCES meetings(MeetingID),

  PRIMARY KEY (TeamID, MeetingID)
);

DROP TABLE IF EXISTS commitments;
CREATE TABLE commitments (
  PersonID  INTEGER NOT NULL REFERENCES people(PersonID),
  MeetingID INTEGER NOT NULL REFERENCES meetings(MeetingID),

  PRIMARY KEY (PersonID, MeetingID)
);

DROP TABLE IF EXISTS signed_in;
CREATE TABLE signed_in (
  PersonID  INTEGER NOT NULL REFERENCES people(PersonID),
  MeetingID INTEGER NOT NULL REFERENCES meetings(MeetingID),
  InTime TIME NOT NULL,

  PRIMARY KEY (PersonID, MeetingID)
);

DROP TABLE IF EXISTS signed_out;
CREATE TABLE signed_out (
  PersonID  INTEGER NOT NULL REFERENCES people(PersonID),
  MeetingID INTEGER NOT NULL REFERENCES meetings(MeetingID),
  outTime TIME NOT NULL,

  PRIMARY KEY (PersonID, MeetingID)
);

DROP TABLE IF EXISTS mentors;
CREATE TABLE mentors(
  PersonID  INTEGER NOT NULL REFERENCES people(PersonID),
  TeamID INTEGER NOT NULL REFERENCES teams(TeamID),

  PRIMARY KEY (PersonID, TeamID)
);

DROP TABLE IF EXISTS students;
CREATE TABLE students (
  PersonID  INTEGER NOT NULL REFERENCES people(PersonID),
  TeamID INTEGER NOT NULL REFERENCES teams(TeamID),

  PRIMARY KEY (PersonID, TeamID)
);

DROP TABLE IF EXISTS parents;
CREATE TABLE parents(
  ParentID INTEGER NOT NULL REFERENCES people(PersonID),
  StudentID INTEGER NOT NULL REFERENCES people(PersonID),

  Relation PARENTRELATION NOT NULL,

  PRIMARY KEY (ParentID, StudentID)
);

--