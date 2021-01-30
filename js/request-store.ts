import { promises as fs, constants } from 'fs';
import { EOL } from 'os';

class NoSuchSubdomain extends Error {

  public subdomain: string;

  constructor ({ subdomain }: { subdomain : string }) {
    super(`No such domain ${subdomain}`);
    this.subdomain = subdomain;
    Error.captureStackTrace(this, NoSuchSubdomain);
  }
}

/**
 * @class RequestsStore Stores and reads requests for a domain.
 */
class RequestsStore {
  private storeDir: string;

  constructor({ storeDir = "data" } : { storeDir?: string } = {}) {
    this.storeDir = storeDir;
  }

  async createStore({ subdomain } : { subdomain : string }) {
    const fileName = this.getFileNameFromSubdomain(subdomain);
  }

  async write({ subdomain, timestamp, content } : { subdomain : string, timestamp : string, content : string }): Promise<void> {
    const fileName = this.getFileNameFromSubdomain(subdomain);
    const contentToBase64 = Buffer.from(content).toString('base64')
    const toAppend = "-----" + timestamp + EOL + contentToBase64 + EOL;

    let file = undefined;

    try {
      file = await fs.open(fileName, constants.O_APPEND);
      await file.write(toAppend, 0, 'utf-8');
    } catch (e) {
      if (e.code === 'ENOENT') {
        throw new NoSuchSubdomain({ subdomain });
      } else {
        throw e;
      }
    } finally {
      file?.close();
    }

  }

  async read({
    subdomain
  } : { subdomain: string }): Promise<string> {
    return await fs.readFile(this.getFileNameFromSubdomain(subdomain), 'utf-8');
  }

  getFileNameFromSubdomain(subdomain: string) {
    return `${this.storeDir}/${subdomain}`;
  }
}

export {
  RequestsStore,
  NoSuchSubdomain
};
