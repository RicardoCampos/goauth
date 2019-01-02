const request = require('request');

module.exports = {
   waitForAppToBecomeAvailable: () => {
    return new Promise((resolve, reject) => {
      const endpoint = (process.env.TEST_ENDPOINT || 'http://localhost:8080');
      const maxWaits = 50;
      let wait = 0;
      const interval = setInterval(() => {
        wait++;
        request(`${endpoint}/metrics`, (err, res) => {
          if (!err && res.statusCode === 200) {
            clearInterval(interval);
            return resolve();
          }
  
          if (wait > maxWaits) {
            clearInterval(interval);
            console.log(`Giving up waiting for ${endpoint}.`);
            reject();
          }
        });
      }, 100);
    });
  }
}