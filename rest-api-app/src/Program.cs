using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Hosting;

namespace rest_api_app
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run(); // Build and run the host
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureWebHostDefaults(webBuilder =>
                {
                    webBuilder.UseStartup<Startup>(); // Specify the startup class
                    webBuilder.UseUrls("http://*:8080"); // Set the application URL to port 8080
                });
    }
}