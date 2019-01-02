const expect = require('chai').expect;
const request = require('superagent');

const { waitForAppToBecomeAvailable } = require('./waitForAppToBecomeAvailable');

before(async () => {
    await waitForAppToBecomeAvailable();
});

describe('The /metrics endpoint',  () => {
    const endpoint = (process.env.TEST_ENDPOINT || 'http://localhost:8080') + '/metrics';
    it('will return 200',  (done) => {
        request.get(endpoint)
            .end((err, response) => {
                expect(err).to.be.null;
                expect(response.status).to.equal(200);
                done();
            });        
    });
    it('will return prometheus metrics',  (done) => {
        request.get(endpoint)
            .end((err, response) => {
                expect(err).to.be.null;
                expect(response.text.startsWith('# HELP go_gc_duration_seconds A summary of the GC invocation durations.')).to.be.true;
                done();
            });        
    });
});