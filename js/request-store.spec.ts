import { RequestsStore, NoSuchSubdomainError } from './request-store';
import { promises as fs } from 'fs';
import chai from 'chai';
import chaiAsPromised from 'chai-as-promised';
import os from 'os';
import path from 'path';

chai.use(chaiAsPromised);
const expect = chai.expect;

describe('RequestsStore', function() {

  const subdomain = 'my-subdomain';

  let dir: string | undefined;
  let store: RequestsStore | undefined;
  beforeEach(async function() {
    dir = await fs.mkdtemp(path.join(os.tmpdir(), 'foo-'));

    store = new RequestsStore({
      storeDir: dir
    });

    if (dir === undefined) {
      throw "Error setting up tests - couldn't create temp directory";
    }
    if (store === undefined) {
      throw "Error setting up tests - couldn't create store";
    }
  });

  afterEach(async function() {
    await fs.rmdir(dir!, { recursive: true });
  });

  describe('#write()', function() {
    it("write doesn't work if subdomain not yet created", async function() {

      await expect(store!.write({
          subdomain: subdomain,
          timestamp: '2020-01-01',
          content: 'my content'
      })).to.be.rejectedWith(NoSuchSubdomainError);

    });
    it("write works if subdomain not yet created", async function() {

      await store!.createStore({ subdomain });

      await expect(store!.write({
          subdomain,
          timestamp: '2020-01-01',
          content: 'my content'
      })).to.not.be.rejected;

    });
    it("write puts out the correct content", async function() {

      await store!.createStore({ subdomain });

      await store!.write({
          subdomain,
          timestamp: '2020-01-01',
          content: 'my content1'
      });
      await store!.write({
          subdomain,
          timestamp: '2020-01-02',
          content: 'my content2'
      });
      await store!.write({
          subdomain,
          timestamp: '2020-01-03',
          content: 'my content3'
      });
      await store!.write({
          subdomain,
          timestamp: '2020-01-04',
          content: 'my content4'
      });
      await store!.write({
          subdomain,
          timestamp: '2020-01-05',
          content: 'my content5'
      });


      expect(await store!.read({subdomain}))
        .to.equal(`-----2020-01-01
bXkgY29udGVudDE=
-----2020-01-02
bXkgY29udGVudDI=
-----2020-01-03
bXkgY29udGVudDM=
-----2020-01-04
bXkgY29udGVudDQ=
-----2020-01-05
bXkgY29udGVudDU=
`);

    });
  });
});
