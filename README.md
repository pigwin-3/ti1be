# ti1be
## Endpoints:
(None of witch are created yet lel)
(diz just a plan for now)
## /journey
### ./get?...
What to return:
```sql
SELECT * FROM public.estimatedvehiclejourney ORDER BY id DESC
```
Peramiters:
- limit
  - default 50
  - ```sql
    LIMIT 50
    ```
  - max 1k (remember to make easy to change
- vehicle_ref
  - ```sql
    Where estimatedvehiclejourney.vehicleref = '1881'
    ```
  - Should also support mutliple lines
  - ```sql
    WHERE estimatedVehicleJourney.vehicleRef IN ('1881', 'SL180000000000479', 'osv')
    ```
- data_source
  - ```sql
    WHERE estimatedVehicleJourney.datasource = 'RUT'
    ```
  - Should also support mutliple lines
  - ```sql
    WHERE estimatedVehicleJourney.datasource IN ('RUT', 'VYG', 'osv')
    ```
- line_ref
  - ```sql
    WHERE estimatedVehicleJourney.lineref = 'RUT:Line:15'
    ```
  - Should also support mutliple lines
  - ```sql
    WHERE estimatedVehicleJourney.lineref IN ('RUT:Line:15', 'TEL:Line:8609', 'osv')
    ```
- after id
  - for pagination
  - ```sql
    id < 1000000
    ```
    - When in use i shal use the last id in previos result

### ./calls?id=...
If not int give 404
Peramiters:
- id
  - ```sql
    estimatedvehiclejourney WHERE id = 943730
    ```
    and
    ```sql
    calls WHERE estimatedvehiclejourney = 943730
    ```
- limit
  - default 200
  - ```sql
    LIMIT 200
    ```
  - max 1k (remember to make easy to change

first get:
```sql
SELECT * FROM estimatedvehiclejourney WHERE id = 943730
```
then get:
```sql
SELECT *
FROM calls
WHERE estimatedvehiclejourney = 943730
ORDER BY "order" ASC
LIMIT 200
```

## /calls
### ./get?id=...
```sql
SELECT * FROM calls WHERE id = 25148238
```

