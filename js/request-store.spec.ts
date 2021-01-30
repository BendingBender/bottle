import { RequestsStore, NoSuchSubdomain } from './request-store';
import { promises as fs } from 'fs';
import chaiAsPromised from 'chai-as-promised';
import os from 'os';
import path from 'path';

describe('RequestsStore', function() {
  let dir: string | undefined = undefined;
  beforeEach(async function() {
    dir = await fs.mkdtemp(path.join(os.tmpdir(), 'foo-'));

    console.log(dir);
  });
  afterEach(async function() {
    if (dir === undefined) {
      throw "Error setting up tests - can't tear down";
    }
    await fs.rmdir(dir);
  });

  describe('#write()', function() {
    it("write doesn't work if subdomain not yet created", async function() {
      if (dir === undefined) {
        throw "Error setting up tests";
      }
      const store = new RequestsStore({
        storeDir: dir
      });

      expect(store.write({
          subdomain: 'my-subdomain',
          timestamp: '2020-01-01',
          content: 'my content'
      })).to.be.rejectedWith(NoSuchSubdomain);

    });
  });
});
