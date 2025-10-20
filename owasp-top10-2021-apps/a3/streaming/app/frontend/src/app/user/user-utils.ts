export class UserUtil {
  static get user() {
    let username = localStorage.getItem('username');

    if(!username) {
      const array = new Uint32Array(1);
      window.crypto.getRandomValues(array);
      let number = array[0] % 100;
      username = `anonymous${number}`;
      localStorage.setItem('username', username);
    }

    return username;
  } 
}