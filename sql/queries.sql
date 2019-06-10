-- Get all teams
SELECT
  compeition,
  number,
  name
FROM teams;

-- Get all People
SELECT
  firstname,
  lastname
FROM people
ORDER BY lastname, firstname;

-- All students and teams
SELECT
  firstname,
  lastname,
  t.compeition,
  t.number,
  t.name
FROM people
  INNER JOIN students s2 on people.personid = s2.personid
  INNER JOIN teams t on s2.teamid = t.teamid
ORDER BY lastname, firstname;

-- Get all students of FRC 3081
SELECT
  people.personid,
  people.checkinid,
  people.firstname,
  people.lastname,
  people.email,
  people.schoolemail,
  people.schoolid
FROM people
  INNER JOIN students s2 on people.personid = s2.personid
  INNER JOIN teams t on s2.teamid = t.teamid
WHERE t.compeition = 'FRC' AND t.number = 3081
ORDER BY lastname, firstname;

-- Mentors
SELECT
  people.personid,
  people.checkinid,
  people.firstname,
  people.lastname,
  people.email,
  people.schoolemail,
  people.schoolid
FROM people
  INNER JOIN mentors m2 on people.personid = m2.personid
  INNER JOIN teams t on m2.teamid = t.teamid
WHERE t.compeition = 'FRC' AND t.number = 3081
ORDER BY lastname, firstname;

-- Get all students of FRC 3081 with parents
-- Get all parents of students of FRC 3081
SELECT
  people.firstname,
  people.lastname,
  student.firstname as student_first,
  student.lastname  as student_last
FROM people
  INNER JOIN parents p on people.personid = p.parentid
  INNER JOIN students s on p.studentid = s.personid
  INNER JOIN people student on student.personid = s.personid
  INNER JOIN teams t on s.teamid = t.teamid
WHERE t.compeition = 'FRC' AND t.number = 3081
ORDER BY lastname, firstname;

SELECT
  student.personid,
  student.checkinid,
  student.firstname,
  student.lastname,
  student.email,
  student.phone,
  student.schoolemail,
  student.schoolid
FROM people as parent
  INNER JOIN parents p on parent.personid = p.parentid
  INNER JOIN students s on p.studentid = s.personid
  INNER JOIN people student on student.personid = s.personid
WHERE parent.personid = 5
ORDER BY student.lastname, student.firstname;

SELECT
  p.relation,
  parent.personid,
  parent.checkinid,
  parent.firstname,
  parent.lastname,
  parent.email,
  parent.phone,
  parent.schoolemail,
  parent.schoolid
FROM people as student
  INNER JOIN students s on student.personid = s.personid
  INNER JOIN parents p on s.personid = p.studentid
  INNER JOIN people parent on parent.personid = p.parentid
WHERE student.personid = 4
ORDER BY parent.lastname, parent.firstname;

SELECT
  parentid,
  studentid,
  relation
FROM parents
WHERE studentid = 4;

-- Get all mentors of FRC 3081
SELECT
  firstname,
  lastname
FROM people
  INNER JOIN mentors m2 on people.personid = m2.personid
  INNER JOIN teams t on m2.teamid = t.teamid
WHERE t.compeition = 'FRC' AND t.number = 3081
ORDER BY lastname, firstname;

-- Teams paul mentors
SELECT
  teams.teamid,
  compeition,
  number,
  name
FROM teams
  INNER JOIN mentors m2 on teams.teamid = m2.teamid
WHERE m2.personid = 0;

-- Person
SELECT
  personid,
  checkinid,
  firstname,
  lastname,
  email,
  phone,
  schoolemail,
  schoolid
FROM people
WHERE personid = 0;

-- Teams that Luke R is on
-- SELECT c


-- Get all mentors of FTC 15530
SELECT
  firstname,
  lastname
FROM people
  INNER JOIN mentors m2 on people.personid = m2.personid
  INNER JOIN teams t on m2.teamid = t.teamid
WHERE t.compeition = 'FTC' AND t.number = 15530
ORDER BY lastname, firstname;

-- Meetings
SELECT *
FROM meetings
ORDER BY date, starttime;

-- 3081 Meetings
SELECT
  meeting.meetingid,
  meeting.date,
  meeting.starttime,
  meeting.endtime,
  meeting.location
FROM meetings as meeting
  INNER JOIN team_meetings m2 on meeting.meetingid = m2.meetingid
  INNER JOIN teams t on m2.teamid = t.teamid
WHERE t.teamid = 1
ORDER BY date, starttime;

SELECT
  teamid,
  compeition,
  number,
  name
FROM teams
WHERE teamid = 0


-- Teams at meeting
SELECT
  t.teamid,
  compeition,
  number,
  name
FROM meetings
  INNER JOIN team_meetings m2 on meetings.meetingid = m2.meetingid
  INNER JOIN teams t on m2.teamid = t.teamid
WHERE m2.meetingid = 0;


SELECT
  meetingid,
  date,
  starttime,
  endtime,
  location
FROM meetings
WHERE meetingid = 0;

-- Get People committed for meeting
SELECT
  person.personid,
  person.checkinid,
  person.firstname,
  person.lastname,
  person.email,
  person.phone,
  person.schoolemail,
  person.schoolid
FROM people as person
  INNER JOIN commitments c2 on person.personid = c2.personid
WHERE c2.meetingid = 0;

-- Get SignedIn at meeting
SELECT
  personid,
  meetingid,
  intime
FROM signed_in
WHERE meetingid = 0