
const RE = /^\w[\w-]*$/;
export const valid = (subdomain: any) => RE.test(subdomain)
