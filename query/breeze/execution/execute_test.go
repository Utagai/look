package execution_test

import "testing"

/*
* This package tests execution, but does so in an integration style. It generally
* works by taking an input query, as a string, and passing it through the parser
* and then executing the AST that it creates. In other words, this is not testing
* execution in isolation.
*
* The downside there is of course that failures in the parser may manifest here as
* well. Ideally however, cases like this are easy enough to determine as being a
* parser issue and don't add extra time trying to determine where it is coming
* from. Plus, always starting debugging at the parser layer is a good rule of
* thumb that can mitigate this problem.
*
* The reason for doing things this way is two-fold:
*   * It is literally easier to write tests this way instead of trying to hand-write
*     the AST for each test case.
*   * It is immune to _good_ or _neutral_ changes to the AST/parser behavior causing
*     these tests to be re-written. They should 'just work'. This is the flip-side
*     of the issue described above, where good or neutral changes that should not
*     cause failures or issues can cause development friction.
**/

func TestQueryExecution(t *testing.T) {}
