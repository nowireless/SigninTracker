-- FTC 15530
-- FRC 3081
INSERT INTO teams(teamid, compeition, number, name)
    VALUES
      (0, 'FTC', 15530, 'Aquila Bots'),
      (1, 'FRC', 3081, 'RoboEagles');

-- Paul Lindemann
--  Mentor of FRC 3081
--  Mentor of FRC 15540
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (0, '000', 'Paul', 'Lindemann');
INSERT INTO mentors(personid, teamid) VALUES (0, 0);
INSERT INTO mentors(personid, teamid) VALUES (0, 1);


-- Ryan Sjostrand
--  Mentor of FRC 3081
--  Mentor of FTC 15530
INSERT INTO people(personid, checkinid, firstname, lastname)
    values (1, '333', 'Ryan', 'Sjostrand');
INSERT INTO mentors(personid, teamid) VALUES (1, 0);
INSERT INTO mentors(personid, teamid) VALUES (1, 1);

-- Sharon Rauenhorst
--  Mentor of FRC 3081
--  Mentor of FTC 15530
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (9, '999', 'Sharon', 'Rauenhorst');
INSERT INTO mentors(personid, teamid) VALUES (9, 0);
INSERT INTO mentors(personid, teamid) VALUES (9, 1);

-- Steve Peterson
--- Mentor of FRC 3081
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (2, '222', 'Steve', 'Peterson');
INSERT INTO mentors(personid, teamid) VALUES (2, 1);

-- Samantha Tranowski
--  Mentor of FRC 3081
--  Mentor of FTC 15530
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (3, '3333', 'Samantha', 'Tarnowski');
INSERT INTO mentors(personid, teamid) values (3, 0);
INSERT INTO mentors(personid, teamid) values (3, 1);

-- Trey Mathieu
--  Student on FRC 3081
--  Mentor on FTC 15530
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (4, '444', 'Trey', 'Mathieu');
INSERT INTO students(personid, teamid) VALUES (4, 1);
INSERT INTO mentors(personid, teamid)  VALUES (4, 0);

-- Rick Mathieu
--  Mentor of FRC 3081
--  Mentor of FTC 15530
--  Parent of Trey
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (5, '555', 'Rick', 'Mathieu');
INSERT INTO mentors(personid, teamid) VALUES (5, 0);
INSERT INTO mentors(personid, teamid) VALUES (5, 1);
INSERT INTO parents(parentid, studentid, relation) VALUES (5, 4, 'Father');


-- Luke Hill
--  Student on FRC 3081
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (6, '666', 'Luke', 'Hill');
INSERT INTO students(personid, teamid) VALUES (6, 1);

-- Jonathan Hill
--  Parent of Luke Hill
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (7, '777', 'Jonathan', 'Hill');
INSERT INTO parents(parentid, studentid, relation) VALUES (7, 6, 'Father');

-- Luke Roberts
--  Student on FRC 3081
--  Student on FTC 15530
INSERT INTO people(personid, checkinid, firstname, lastname)
    VALUES (8, '888', 'Luke', 'Roberts');
INSERT INTO students(personid, teamid) VALUES (8, 0);
INSERT INTO students(personid, teamid) VALUES (8, 1);


-- Times in ISO 8601

-- Create meeting
INSERT INTO meetings(meetingid, date, starttime, endtime, location, kind)
    VALUES (0, '6-8-2019', '10:00:00 CST', '12:00:00 CST', 'Kennedy High School', 'General');

INSERT INTO team_meetings(teamid, meetingid) VALUES (1, 0);

-- Sign ins

-- Paul
INSERT INTO signed_in(personid, meetingid, intime)
    VALUES (0, 0, '10:00:00 CST');
INSERT INTO signed_out(personid, meetingid, outtime)
    VALUES (0, 0, '12:00:00 CST');

-- Ryan
INSERT INTO commitments(personid, meetingid)
    VALUES (1, 0);
INSERT INTO signed_in(personid, meetingid, intime)
    VALUES (1, 0, '10:00:00 CST');
INSERT INTO signed_out(personid, meetingid, outtime)
    VALUES (1, 0, '1:30:00 CST');

-- Samantha
INSERT INTO signed_in(personid, meetingid, intime)
    VALUES (3, 0, '10:00:00 CST');
INSERT INTO signed_out(personid, meetingid, outtime)
    VALUES (3, 0, '1:30:00 CST');

-- Trey
INSERT INTO signed_in(personid, meetingid, intime)
    VALUES (4, 0, '10:00:00 CST');
INSERT INTO signed_out(personid, meetingid, outtime)
    VALUES (4, 0, '1:30:00 CST');

--  Luke
INSERT INTO signed_in(personid, meetingid, intime)
    VALUES (6, 0, '10:00:00 CST');
INSERT INTO signed_out(personid, meetingid, outtime)
    VALUES (6, 0, '12:00:00 CST');

-- Sharon
INSERT INTO signed_in(personid, meetingid, intime)
    VALUES (9, 0, '10:00:00 CST');
INSERT INTO signed_out(personid, meetingid, outtime)
    VALUES (9, 0, '12:00:00 CST');