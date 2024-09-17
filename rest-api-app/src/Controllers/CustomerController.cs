using Microsoft.AspNetCore.Mvc;
using System.Collections.Generic;

namespace rest_api_app.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class CustomerController : ControllerBase
    {
        private readonly List<Customer> _customers;

        public CustomerController()
        {
            _customers = new List<Customer>
            {
                new Customer
                {
                    LoginId = "john123",
                    Name = "John Doe",
                    Email = "john.doe@example.com",
                    Address = "123 Main St"
                },
                new Customer
                {
                    LoginId = "jane456",
                    Name = "Jane Smith",
                    Email = "jane.smith@example.com",
                    Address = "456 Elm St"
                }
            };
        }

        [HttpGet("getCustomerInfo")]
        public IActionResult GetCustomerInfo()
        {
            return Ok(_customers);
        }
    }

    public class Customer
    {
        public string LoginId { get; set; }
        public string Name { get; set; }
        public string Email { get; set; }
        public string Address { get; set; }
    }
}