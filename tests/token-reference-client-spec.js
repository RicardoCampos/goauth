const expect = require('chai').expect;
const request = require('superagent');
const validate = require('uuid-validate');

const { waitForAppToBecomeAvailable } = require('./waitForAppToBecomeAvailable');

before(async () => {
    await waitForAppToBecomeAvailable();
});

describe('The /token endpoint with a reference token client',  () => {
    const endpoint = (process.env.TEST_ENDPOINT || 'http://localhost:8080') + '/token';
    it('will return a 200 response',  (done) => {
      request.post(endpoint)
        .auth('foo_reference', 'secret')
        .send('grant_type=client_credentials')
        .send('scope=read')
        .end((err, response) => {
            expect(err).to.be.null;
            expect(response.status).to.equal(200);
            done();
        });        
    });
    it('will return an reference access token',  (done) => {
        request.post(endpoint)
          .auth('foo_reference', 'secret')
          .send('grant_type=client_credentials')
          .send('scope=read')
          .end((err, response) => {
              expect(err).to.be.null;
              const result = JSON.parse(response.text);
              expect(validate(result.access_token,4)).to.be.true;
              done();
          });        
    });
    it('will return a scope',  (done) => {
      request.post(endpoint)
        .auth('foo_reference', 'secret')
        .send('grant_type=client_credentials')
        .send('scope=read')
        .end((err, response) => {
            expect(err).to.be.null;
            const result = JSON.parse(response.text);
            expect(result.scope).to.equal('read');
            done();
        });        
  });
  it('will return an reference token type',  (done) => {
    request.post(endpoint)
      .auth('foo_reference', 'secret')
      .send('grant_type=client_credentials')
      .send('scope=read')
      .end((err, response) => {
          expect(err).to.be.null;
          const result = JSON.parse(response.text);
          expect(result.token_type).to.equal('Reference');
          done();
      });        
  });
});