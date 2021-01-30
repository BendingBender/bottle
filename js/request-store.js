const { fs } = require('fs');
const { EOL } = require('os');

/**
 * @class RequestsStore Stores and reads requests for a domain.
 */
class RequestsStore {

  constructor({
    name
  }) {
    this.fileName = `data/${name}`;
  }

  async write({ timestamp, content }) {
    const contentToBase64 = Buffer.from(content).toString('base64')
    const toAppend = "-----" + timestamp + EOL + contentToBase64 + EOL;

    await fs.appendFile(this.fileName, toAppend, {
      encoding: 'utf8'
    });
  }

  async read() {
    return await fs.readFile(this.fileName);
  }

}

module.exports = {
  RequestsStore
}