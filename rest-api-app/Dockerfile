FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 80

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["rest-api-app.csproj", "."]
RUN dotnet restore "./rest-api-app.csproj"
COPY . .
WORKDIR "/src/"
RUN dotnet build "rest-api-app.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "rest-api-app.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "rest-api-app.dll"]