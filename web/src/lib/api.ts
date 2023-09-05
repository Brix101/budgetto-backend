import axios from "axios";

const BE_BASE_URL = "http://192.168.254.180:5000/api";
export default axios.create({
  baseURL: BE_BASE_URL,
  withCredentials: true,
});
