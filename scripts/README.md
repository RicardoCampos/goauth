You can create the required Postgres database (called "goauth") by running `createPostgresDb.sql` in an admin tool or `psql`.

If you want to insert some test data, run `insertTestClients.sql` in the same way. Of course, don't do this in production.

`initialise_test_db.sh` is used for the test Postgres Dockerfile. This simply runs the two above, ensuring that the commands run on the correct database (there is no `USE Database` in Postgres).