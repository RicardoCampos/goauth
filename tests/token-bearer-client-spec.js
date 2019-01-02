const expect = require('chai').expect;
const request = require('superagent');
const jwt = require('jsonwebtoken');
const { waitForAppToBecomeAvailable } = require('./waitForAppToBecomeAvailable');

before(async () => {
    await waitForAppToBecomeAvailable();
});

describe('The /token endpoint with a bearer token client',  () => {
    const endpoint = (process.env.TEST_ENDPOINT || 'http://localhost:8080') + '/token';
    it('will return a 200 response',  (done) => {
      request.post(endpoint)
        .auth('foo_bearer', 'secret')
        .send('grant_type=client_credentials')
        .send('scope=read')
        .end((err, response) => {
            expect(err).to.be.null;
            expect(response.status).to.equal(200);
            done();
        });        
    });
    it('will return an bearer access token',  (done) => {
        request.post(endpoint)
          .auth('foo_bearer', 'secret')
          .send('grant_type=client_credentials')
          .send('scope=read')
          .end((err, response) => {
              expect(err).to.be.null;
              const result = JSON.parse(response.text);
              const token = jwt.decode(result.access_token, {complete: true});
              expect(token.header.alg).to.equal('RS512');
              expect(token.header.typ).to.equal('JWT');
              done();
          });        
    });
    it('will return a scope',  (done) => {
      request.post(endpoint)
        .auth('foo_bearer', 'secret')
        .send('grant_type=client_credentials')
        .send('scope=read')
        .end((err, response) => {
            expect(err).to.be.null;
            const result = JSON.parse(response.text);
            expect(result.scope).to.equal('read');
            done();
        });        
  });
  it('will return an bearer token type',  (done) => {
    request.post(endpoint)
      .auth('foo_bearer', 'secret')
      .send('grant_type=client_credentials')
      .send('scope=read')
      .end((err, response) => {
          expect(err).to.be.null;
          const result = JSON.parse(response.text);
          expect(result.token_type).to.equal('Bearer');
          done();
      });        
  });
});