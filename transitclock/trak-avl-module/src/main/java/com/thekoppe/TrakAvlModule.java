package com.thekoppe;

import java.io.InputStream;
import java.net.URLConnection;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Date;

import org.json.JSONObject;
import org.json.JSONArray;
import org.json.JSONException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.transitclock.avl.PollUrlAvlModule;
import org.transitclock.config.StringConfigValue;
import org.transitclock.db.structs.AvlReport;
import org.transitclock.modules.Module;

public class TrakAvlModule extends PollUrlAvlModule {
  private static StringConfigValue baseUrl = new StringConfigValue("transitclock.avl.trak.url",
      "https://beartrakapi.thekoppe.com",
      "The base URL of the Trak API to use.");

  private static final DateFormat trakTimeFormat = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ssXXX");

  private static final Logger logger = LoggerFactory.getLogger(TrakAvlModule.class);

  /********************** Member Functions **************************/

  public TrakAvlModule(String agencyId) {
    super(agencyId);
  }

  @Override
  protected String getUrl() {
    return baseUrl.getValue() + "/v1/transit/vehicles";
  }

  @Override
  protected void setRequestHeaders(URLConnection con) {
    con.addRequestProperty("User-Agent",
        "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.0)");
  }

  /**
   * Reads in the JSON data from the InputStream and creates and then
   * processes an AvlReport.
   * 
   * @param in
   * @return Collection of AvlReports
   */
  @Override
  protected Collection<AvlReport> processData(InputStream in) throws Exception {
    String jsonStr = getJsonString(in);

    try {
      JSONArray jsonArray = new JSONArray(jsonStr);

      Collection<AvlReport> avlReportsReadIn = new ArrayList<AvlReport>();

      // process data for each vehicle
      for (int i = 0; i < jsonArray.length(); ++i) {
        JSONObject vehicleData = jsonArray.getJSONObject(i);

        Object idObj = vehicleData.get("id");
        String vehicleId = (idObj instanceof String) ? (String) idObj : String.valueOf(idObj);

        double lat = vehicleData.getDouble("latitude");
        double lon = vehicleData.getDouble("longitude");
        float heading = (float) vehicleData.getInt("heading");
        float speed = (float) vehicleData.getDouble("speed");

        String updatedTimeStr = vehicleData.getString("lastUpdated");
        Date updatedTime = trakTimeFormat.parse(updatedTimeStr);

        AvlReport avlReport = new AvlReport(vehicleId, updatedTime.getTime(), lat, lon, speed, heading, "Trak");
        avlReportsReadIn.add(avlReport);
      }

      return avlReportsReadIn;
    } catch (JSONException e) {
      logger.error("Error parsing JSON. {}. {}",
          e.getMessage(), jsonStr, e);
      return new ArrayList<AvlReport>();
    }
  }

  /**
   * Just for debugging
   */
  public static void main(String[] args) {
    Module.start("org.transitclock.avl.TranslocAvlModule");
  }
}
