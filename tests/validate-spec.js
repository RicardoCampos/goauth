const expect = require('chai').expect;
const request = require('superagent');

const { waitForAppToBecomeAvailable } = require('./waitForAppToBecomeAvailable');

before(async () => {
    await waitForAppToBecomeAvailable();
});

describe('The /validate endpoint',  () => {
    const endpoint = (process.env.TEST_ENDPOINT || 'http://localhost:8080');
    it('with a valid reference token will return a 200 response',  (done) => {
      request.post(endpoint + '/token')
        .auth('foo_reference', 'secret')
        .send('grant_type=client_credentials')
        .send('scope=read')
        .end((err, response) => {
            expect(err).to.be.null;
            expect(response.status).to.equal(200);
            const result = JSON.parse(response.text);
            request.post(endpoint + '/validate')
              .send(`token=${result.access_token}`)
              .end((err, response) => {
                expect(err).to.be.null;
                expect(response.status).to.equal(200);
                done();
              });  
        });        
    });
    it('with an invalid reference token will return a 400 response',  (done) => {
        request.post(endpoint + '/validate')
          .send(`token=notarealtoken`)
          .end((response) => {
              expect(response.status).to.equal(400); 
              done();
          });        
      });
});