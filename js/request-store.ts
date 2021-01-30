import { promises as fs, constants } from 'fs';
import { EOL } from 'os';

/**
 * @class NoSuchSubdomainError is thrown when you try to write a request
 * to an empty subdomain.
 */
export class NoSuchSubdomainError extends Error {

  public subdomain: string;

  constructor ({ subdomain }: { subdomain : string }) {
    super(`No such domain ${subdomain}`);
    this.subdomain = subdomain;
  }
}

/**
 * @class RequestsStore stores and reads requests for a domain.
 */
export class RequestsStore {
  private storeDir: string;

  constructor({ storeDir } : { storeDir: string }) {
    this.storeDir = storeDir;
  }

  async createStore({ subdomain } : { subdomain : string }) {
    const fileName = this.getFileNameFromSubdomain(subdomain);
    await fs.writeFile(fileName, '', 'utf-8');
  }

  async exists({ subdomain } : { subdomain : string }): Promise<boolean> {
    try {
      const fileName = this.getFileNameFromSubdomain(subdomain);
      await fs.access(fileName, constants.F_OK);
      return true;
    } catch {
      return false;
    }
  }

  async write({ subdomain, timestamp, content } : { subdomain : string, timestamp : string, content : string }): Promise<void> {
    const fileName = this.getFileNameFromSubdomain(subdomain);
    const contentToBase64 = Buffer.from(content).toString('base64')
    const toAppend = "-----" + timestamp + EOL + contentToBase64 + EOL;

    let file = undefined;

    try {
      file = await fs.open(fileName, constants.O_WRONLY | constants.O_APPEND);
      await file.write(toAppend, 0, 'utf-8');
    } catch (e) {
      if (e.code === 'ENOENT') {
        throw new NoSuchSubdomainError({ subdomain });
      } else {
        throw e;
      }
    } finally {
      await file?.close();
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
