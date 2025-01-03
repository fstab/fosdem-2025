import http from 'k6/http';

export const options = {
  // run 10 requests in parallel
  vus: 10,
  iterations: 10,
};

export default function () {
  http.get('http://product-service:8080/products?search=telescope');
}
