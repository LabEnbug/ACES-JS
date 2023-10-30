export default function GetUserInfo() {
  return window.localStorage.getItem('userInfo') ? JSON.parse(window.localStorage.getItem('userInfo')): null;
}
