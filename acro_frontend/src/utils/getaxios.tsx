import axios from 'axios';

export default function GetUserInfo() {
  const userinfo = window.localStorage.getItem('userInfo')
    ? JSON.parse(window.localStorage.getItem('userInfo'))
    : null;
  return userinfo
    ? axios.create({
        headers: {
          Authorization: `Bearer ${userinfo.data.token}`,
        },
      })
    : axios.create();
}

// const baxios = axios.create({
//   // baseURL: ,
//   headers: {
//     // Authorization: `Bearer ${token}`,
//   }
// })

// export default baxios;
